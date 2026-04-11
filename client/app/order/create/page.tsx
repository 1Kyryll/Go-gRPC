import { redirect } from "next/navigation";

export default function CreateOrderPage() {
    const addOrderAction = async (formData: FormData) => {
        "use server";

        const response = await fetch("http://localhost:8080/order/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                customer_id: Number(formData.get("customer_id")),
                items: String(formData.get("items")).split(",").map(s => s.trim()),
            }),
        });

        const data = await response.json();
        redirect(`/ticket/${data.order_id}`);
    }

    return (
        <div className="p-8">
            <h1 className="text-3xl font-bold mb-6">Create Order</h1>
            <form action={addOrderAction} className="space-y-4 max-w-md">
                <input
                    type="number"
                    name="customer_id"
                    placeholder="Customer ID"
                    required
                    className="w-full rounded border border-gray-300 px-3 py-2"
                />
                <input
                    type="text"
                    name="items"
                    placeholder="Items (comma separated)"
                    required
                    className="w-full rounded border border-gray-300 px-3 py-2"
                />
                <button
                    type="submit"
                    className="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
                >
                    Create Order
                </button>
            </form>
        </div>
    );
}
