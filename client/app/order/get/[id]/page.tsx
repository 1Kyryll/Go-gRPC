type Order = {
    id: number;
    customer_id: number;
    items: string[];
    status: string;
};

export default async function GetOrderPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = await params;
    const res = await fetch(`http://localhost:8080/order/get/${id}`, { cache: "no-store" });
    const orders: Order[] = await res.json();

    return (
        <div>
            {orders?.map((order: Order) => (
                <div key={order.id}>
                    <h2>Order #{order.id}</h2>
                    <p>Customer ID: {order.customer_id}</p>
                    <p>Items: {order.items.join(", ")}</p>
                    <p>Status: {order.status}</p>
                </div>
            ))}
        </div>
    );
} 

