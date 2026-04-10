"use client";

import { useEffect, useState } from "react";

type Order = {
  id: number;
  customer_id: number;
  items: string[];
  status: string;
};

export default function KitchenPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const eventSource = new EventSource("http://localhost:8081/stream");

    eventSource.onopen = () => {
      setConnected(true);
    };

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.connected) return;
      setOrders((prev) => [data, ...prev]);
    };

    eventSource.onerror = () => {
      setConnected(false);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="p-8">
      <div className="flex items-center gap-3 mb-6">
        <h1 className="text-3xl font-bold">Kitchen Orders</h1>
        <span
          className={`inline-block w-3 h-3 rounded-full ${
            connected ? "bg-green-500" : "bg-red-500"
          }`}
        />
      </div>

      {orders.length === 0 && (
        <p className="text-gray-500">Waiting for new orders...</p>
      )}

      <div className="space-y-4">
        {orders.map((order, index) => (
          <div
            key={`${order.id}-${index}`}
            className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
          >
            <div className="flex justify-between items-center mb-2">
              <span className="font-bold text-lg">Order #{order.id}</span>
              <span className="text-sm text-gray-500">
                Customer {order.customer_id}
              </span>
            </div>
            <ul className="list-disc list-inside text-gray-700">
              {order.items.map((item, i) => (
                <li key={i}>{item}</li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </div>
  );
}
