import { gql } from "@/lib/graphql";
import Link from "next/link";
import OrderStatusLive from "@/components/OrderStatusLive";
import Navbar from "@/components/Navbar";

type OrderDetail = {
  id: string;
  status: string;
  totalPrice: number;
  createdAt: string;
  updatedAt: string;
  user: { id: string; username: string; email: string } | null;
  items: {
    id: string;
    menuItem: { name: string; price: number } | null;
    quantity: number;
    specialInstructions: string;
    subtotal: number;
  }[];
  ticket: { id: string; status: string; createdAt: string; updatedAt: string } | null;
};

const statusColors: Record<string, string> = {
  PENDING: "bg-gray-100 text-gray-700",
  CONFIRMED: "bg-blue-100 text-blue-700",
  COOKING: "bg-orange-100 text-orange-700",
  COMPLETED: "bg-green-100 text-green-700",
  CANCELLED: "bg-red-100 text-red-700",
};

const ticketStatusColors: Record<string, string> = {
  OPEN: "bg-amber-100 text-amber-700",
  COOKING: "bg-orange-100 text-orange-700",
  DONE: "bg-green-100 text-green-700",
};

export default async function OrderDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let order: OrderDetail | null = null;
  let error: string | null = null;

  try {
    const data = await gql<{ order: OrderDetail }>(
      `query GetOrder($id: ID!) {
        order(id: $id) {
          id
          status
          totalPrice
          createdAt
          updatedAt
          user { id, username, email }
          items {
            id
            menuItem { name, price }
            quantity
            specialInstructions
            subtotal
          }
          ticket { id, status, createdAt, updatedAt }
        }
      }`,
      { id }
    );
    order = data.order;
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load order";
  }

  if (error || !order) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold mb-2">Order not found</h1>
        <p className="text-gray-500 mb-4">
          {error || `Order #${id} doesn't exist.`}
        </p>
        <Link
          href="/orders"
          className="text-emerald-600 font-medium hover:underline"
        >
          Back to Orders
        </Link>
      </div>
    );
  }

  return (
    <>
      <Navbar />
      <div className="max-w-2xl mx-auto px-4 sm:px-6 py-8">
        <Link
          href="/orders"
          className="text-sm text-gray-500 hover:text-gray-700 mb-4 inline-flex items-center gap-1"
        >
          <svg
            className="w-4 h-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          All Orders
        </Link>

        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Order #{order.id}</h1>
          <OrderStatusLive
            orderId={order.id}
            initialStatus={order.status}
            statusColors={statusColors}
          />
        </div>

        <div className="space-y-4">
          {/* User info */}
          {order.user && (
            <div className="bg-white rounded-lg border border-gray-200 p-4 shadow-sm">
              <h2 className="text-sm font-medium text-gray-500 mb-2">
                User
              </h2>
              <p className="font-medium">{order.user.username}</p>
              <p className="text-sm text-gray-500">{order.user.email}</p>
            </div>
          )}

          {/* Order items */}
          <div className="bg-white rounded-lg border border-gray-200 shadow-sm">
            <h2 className="text-sm font-medium text-gray-500 p-4 pb-2">
              Items
            </h2>
            <div className="divide-y divide-gray-100">
              {order.items.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between px-4 py-3"
                >
                  <div>
                    <p className="font-medium">
                      {item.menuItem?.name || "Unknown item"}
                    </p>
                    {item.specialInstructions && (
                      <p className="text-xs text-gray-400 italic mt-0.5">
                        {item.specialInstructions}
                      </p>
                    )}
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-500">
                      {item.quantity} x $
                      {item.menuItem?.price?.toFixed(2) ?? "—"}
                    </p>
                    <p className="font-medium">${item.subtotal.toFixed(2)}</p>
                  </div>
                </div>
              ))}
            </div>
            <div className="flex justify-between px-4 py-3 bg-gray-50 rounded-b-lg border-t border-gray-100">
              <span className="font-bold">Total</span>
              <span className="font-bold">${order.totalPrice.toFixed(2)}</span>
            </div>
          </div>

          {/* Ticket */}
          {order.ticket && (
            <div className="bg-white rounded-lg border border-gray-200 p-4 shadow-sm">
              <div className="flex items-center justify-between mb-2">
                <h2 className="text-sm font-medium text-gray-500">
                  Kitchen Ticket
                </h2>
                <Link
                  href={`/ticket/${order.id}`}
                  className="text-xs text-emerald-600 hover:underline"
                >
                  View details
                </Link>
              </div>
              <div className="flex items-center justify-between">
                <span className="font-medium">Ticket #{order.ticket.id}</span>
                <span
                  className={`text-xs font-medium px-2.5 py-1 rounded-full ${
                    ticketStatusColors[order.ticket.status] ||
                    "bg-gray-100 text-gray-700"
                  }`}
                >
                  {order.ticket.status}
                </span>
              </div>
            </div>
          )}

          {/* Timestamps */}
          <p className="text-xs text-gray-400 text-center">
            Created{" "}
            {new Date(order.createdAt).toLocaleString(undefined, {
              dateStyle: "medium",
              timeStyle: "short",
            })}
          </p>
        </div>
      </div>
    </>
  );
}
