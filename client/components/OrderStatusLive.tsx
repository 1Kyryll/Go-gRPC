"use client";

import { useEffect, useState } from "react";

const GRAPHQL_WS_URL =
  process.env.NEXT_PUBLIC_GRAPHQL_WS_URL ||
  (typeof window !== "undefined"
    ? `ws://${window.location.hostname}:8082/graphql`
    : "ws://localhost:8082/graphql");

export default function OrderStatusLive({
  orderId,
  initialStatus,
  statusColors,
}: {
  orderId: string;
  initialStatus: string;
  statusColors: Record<string, string>;
}) {
  const [status, setStatus] = useState(initialStatus);
  const [flash, setFlash] = useState(false);

  useEffect(() => {
    let ws: WebSocket | null = null;

    try {
      ws = new WebSocket(GRAPHQL_WS_URL, "graphql-transport-ws");

      ws.onopen = () => {
        ws?.send(JSON.stringify({ type: "connection_init" }));
      };

      ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);

        if (msg.type === "connection_ack") {
          ws?.send(
            JSON.stringify({
              id: "1",
              type: "subscribe",
              payload: {
                query: `subscription OrderStatusChanged($orderId: ID) {
                  orderStatusChanged(orderId: $orderId) {
                    id
                    status
                  }
                }`,
                variables: { orderId },
              },
            })
          );
        }

        if (msg.type === "next" && msg.payload?.data?.orderStatusChanged) {
          const newStatus = msg.payload.data.orderStatusChanged.status;
          setStatus(newStatus);
          setFlash(true);
          setTimeout(() => setFlash(false), 1500);
        }
      };
    } catch {
      // WebSocket not available, stay with initial status
    }

    return () => {
      if (ws && ws.readyState <= WebSocket.OPEN) {
        ws.close();
      }
    };
  }, [orderId]);

  return (
    <span
      className={`text-xs font-medium px-2.5 py-1 rounded-full transition-all duration-300 ${
        statusColors[status] || "bg-gray-100 text-gray-700"
      } ${flash ? "scale-110 ring-2 ring-indigo-300" : ""}`}
    >
      {status}
    </span>
  );
}
