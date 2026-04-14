const GRAPHQL_URL = process.env.NEXT_PUBLIC_GRAPHQL_URL || "http://localhost:8082/graphql";

export async function gql<T>(
    query: string, 
    variables?: Record<string, unknown>
): Promise<T> {
    const res = await fetch(GRAPHQL_URL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ query, variables }),
    });

    const json = await res.json();
    
    if (json.errors) {
        throw new Error(json.errors[0].message || "GraphQL error");
    }

    return json.data as T;
}