import { cookies } from "next/headers";

const GRAPHQL_URL = process.env.NEXT_PUBLIC_GRAPHQL_URL || "http://localhost:8082/graphql";

export async function gql<T>(
    query: string,
    variables?: Record<string, unknown>
): Promise<T> {
    const headers: Record<string, string> = {
        "Content-Type": "application/json",
    };

    try {
        const cookieStore = await cookies();
        const token = cookieStore.get("token")?.value;
        if (token) {
            headers["Authorization"] = `Bearer ${token}`;
        }
    } catch {
        // cookies() throws in non-server contexts — skip auth header
    }

    const res = await fetch(GRAPHQL_URL, {
        method: "POST",
        headers,
        body: JSON.stringify({ query, variables }),
    });

    const json = await res.json();

    if (json.errors) {
        throw new Error(json.errors[0].message || "GraphQL error");
    }

    return json.data as T;
}