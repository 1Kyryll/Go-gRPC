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
        
        console.log(await response.json());
    }    
    
    return (
        <div>
            <form action={addOrderAction}>
                <input type="number" name="customer_id" placeholder="Customer ID" required />
                <input type="text" name="items" placeholder="Items (comma separated)" required />
                <button type="submit">Create Order</button>
            </form>
        </div>
    );
} 