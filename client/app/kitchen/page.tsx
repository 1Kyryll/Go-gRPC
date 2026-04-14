"use client";

import { useEffect, useState, useRef } from "react";

type Order = {
  id: number;
  customer_id: number;
  items?: string[];
  status: string;
};

const KITCHEN_STREAM_URL =
  process.env.NEXT_PUBLIC_KITCHEN_STREAM_URL || "http://localhost:8081/stream";
const KITCHEN_API_URL =
  process.env.NEXT_PUBLIC_KITCHEN_API_URL || "http://localhost:8081";

export default function KitchenPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [connected, setConnected] = useState(false);
  const [completing, setCompleting] = useState<number | null>(null);
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const eventSource = new EventSource(KITCHEN_STREAM_URL);

    eventSource.onopen = () => setConnected(true);

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.connected) return;
      setOrders((prev) => [data, ...prev]);
    };

    eventSource.onerror = () => setConnected(false);

    return () => eventSource.close();
  }, []);

  async function handleDone(orderId: number) {
    setCompleting(orderId);
    try {
      const res = await fetch(`${KITCHEN_API_URL}/order/${orderId}/done`, {
        method: "POST",
      });
      if (res.ok) {
        setOrders((prev) => prev.filter((o) => o.id !== orderId));
      }
    } finally {
      setCompleting(null);
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold">Kitchen</h1>
          <span
            className={`inline-flex items-center gap-1.5 text-xs font-medium px-2.5 py-1 rounded-full ${
              connected
                ? "bg-green-100 text-green-700"
                : "bg-red-100 text-red-700"
            }`}
          >
            <span
              className={`w-1.5 h-1.5 rounded-full ${
                connected ? "bg-green-500 animate-pulse" : "bg-red-500"
              }`}
            />
            {connected ? "Live" : "Disconnected"}
          </span>
        </div>
        <span className="text-sm text-gray-500">
          {orders.length} pending{orders.length === 1 ? "" : " orders"}
        </span>
      </div>

      {orders.length === 0 && (
        <div className="text-center py-16">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg
              className="w-8 h-8 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <p className="text-gray-500">Waiting for new orders...</p>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4" ref={listRef}>
        {orders.map((order, index) => (
          <div
            key={`${order.id}-${index}`}
            className="rounded-lg border border-gray-200 bg-white shadow-sm overflow-hidden animate-[fadeSlideIn_0.3s_ease-out]"
          >
            <div className="flex items-center justify-between px-4 py-3 bg-gray-50 border-b border-gray-100">
              <span className="font-bold">Order #{order.id}</span>
              <span className="text-xs text-gray-500 bg-gray-200 px-2 py-0.5 rounded-full">
                Customer {order.customer_id}
              </span>
            </div>
            <div className="p-4">
              {order.items && order.items.length > 0 ? (
                <ul className="space-y-1 mb-4">
                  {order.items.map((item, i) => (
                    <li
                      key={i}
                      className="text-sm text-gray-700 flex items-center gap-2"
                    >
                      <span className="w-1.5 h-1.5 bg-indigo-400 rounded-full shrink-0" />
                      {item}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-sm text-gray-400 italic mb-4">
                  Order #{order.id} — {order.status}
                </p>
              )}
              <button
                onClick={() => handleDone(order.id)}
                disabled={completing === order.id}
                className="w-full rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {completing === order.id ? "Completing..." : "Mark Done"}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
