"use client";

import { useEffect, useState } from "react";

type Order = {  
    id: string;
    customerName: string;
    items: string[];
    status: string;
};

export default function GetOrderPage() {
    const [orders, setOrders] = useState<Order[]>([]);

    useEffect(() => {
        fetch("http://localhost:8080/order/get")
            .then(response => response.json())
            .then(data => setOrders(data))
            .catch(error => console.error("Error fetching orders:", error));
    }, []);
    
    return (
        <div>
            {orders?.map((order) => (
                <div key={order.id}>
                    <h2>Order ID: {order.id}</h2>
                    <p>Customer Name: {order.customerName}</p>
                    <p>Items: {order.items.join(", ")}</p>
                    <p>Status: {order.status}</p>
                </div>
            ))}
        </div>
    );
} 