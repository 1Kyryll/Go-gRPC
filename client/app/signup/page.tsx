export default function SignupPage() {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
            <h1 className="text-3xl font-bold mb-6 text-center">Sign up</h1>
            <form className="bg-white p-6 rounded-lg shadow-md w-full max-w-sm flex flex-col gap-4">
                <div>
                    <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">
                        Username
                    </label>
                    <input
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 border-none leading-tight focus:outline-none focus:shadow-outline"
                        id="username"
                        type="text"
                        placeholder="e. g. John Doe"
                    />
                </div>
                <div>
                    <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">
                        Email
                    </label>
                    <input
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 border-none leading-tight focus:outline-none focus:shadow-outline"
                        id="email"
                        type="email"
                        placeholder="e. g. john.doe@example.com"
                    />
                </div>
                <div>
                    <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="phone">
                        Phone
                    </label>
                    <input
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 border-none leading-tight focus:outline-none focus:shadow-outline"
                        id="phone"
                        type="tel"
                        placeholder="e. g. +1 234 567 890 (optional)"
                    />
                </div>
                <div>
                    <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">
                        Password
                    </label>
                    <input
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 border-none mb-3 leading-tight focus:outline-none focus:shadow-outline"
                        id="password"
                        type="password"
                        placeholder="********"
                    />
                </div>
                <button
                    className="bg-emerald-600 hover:bg-emerald-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                    type="submit"
                >
                    Sign up
                </button>
            </form>
        </div>
    );
}