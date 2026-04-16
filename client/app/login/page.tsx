"use client";

import Link from "next/link";
import { useActionState } from "react";
import { loginAction } from "@/actions/auth";

type FormState = { success: boolean; error?: string };

export default function LoginPage() {
    const [state, formAction, pending] = useActionState(
        async (_prev: FormState, formData: FormData): Promise<FormState> => {
            const result = await loginAction(formData);
            return result;
        },
        { success: true }
    );

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            <h1 className="text-3xl font-bold mb-6 text-center">Login</h1>
            <form action={formAction} className="bg-white p-6 rounded-lg shadow-md w-full max-w-sm flex flex-col gap-4">
                {state.error && (
                    <div className="bg-red-50 text-red-600 text-sm p-3 rounded-md">
                        {state.error}
                    </div>
                )}
                <div>
                    <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">
                        Username
                    </label>
                    <input
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 border-none leading-tight focus:outline-none focus:shadow-outline"
                        id="username"
                        name="username"
                        type="text"
                        placeholder="e.g. johndoe"
                        required
                    />
                </div>
                <div>
                    <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">
                        Password
                    </label>
                    <input
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 border-none mb-3 leading-tight focus:outline-none focus:shadow-outline"
                        id="password"
                        name="password"
                        type="password"
                        placeholder="********"
                        required
                    />
                </div>
                <button
                    className="bg-emerald-600 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50"
                    type="submit"
                    disabled={pending}
                >
                    {pending ? "Logging in..." : "Login"}
                </button>
                <p className="text-sm text-center text-gray-600">
                    Don&apos;t have an account?{" "}
                    <Link href="/signup" className="text-emerald-600 hover:underline">
                        Sign up
                    </Link>
                </p>
            </form>
        </div>
    );
}
