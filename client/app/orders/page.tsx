import { gql } from "@/lib/graphql";
import Link from "next/link";
import OrderStatusFilter from "@/components/OrderStatusFilter";
import Navbar from "@/components/Navbar";

type Order = {
  id: string;
  status: string;
  totalPrice: number;
  createdAt: string;
  customer: { id: string; name: string } | null;
};

type Props = {
  searchParams: Promise<{ status?: string }>;
};

const statusColors: Record<string, string> = {
  PENDING: "bg-gray-100 text-gray-700",
  CONFIRMED: "bg-blue-100 text-blue-700",
  COOKING: "bg-orange-100 text-orange-700",
  COMPLETED: "bg-green-100 text-green-700",
  CANCELLED: "bg-red-100 text-red-700",
};

export default async function OrdersPage({ searchParams }: Props) {
  const { status } = await searchParams;

  const validStatuses = [
    "PENDING",
    "CONFIRMED",
    "COOKING",
    "COMPLETED",
    "CANCELLED",
  ];
  const selectedStatus =
    status && validStatuses.includes(status) ? status : null;

  const data = await gql<{
    orders: {
      edges: { node: Order }[];
      totalCount: number;
    };
  }>(
    `query GetOrders($status: OrderStatus, $first: Int) {
      orders(status: $status, first: $first) {
        edges {
          node {
            id
            status
            totalPrice
            createdAt
            customer { id, name }
          }
        }
        totalCount
      }
    }`,
    { status: selectedStatus, first: 50 }
  );

  const orders = data.orders.edges.map((e) => e.node);

  return (
    <>
      <Navbar />
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6">
          <h1 className="text-2xl font-bold mb-1">Orders</h1>
          <p className="text-sm text-gray-500">
            {data.orders.totalCount} orders
          </p>
        </div>

        <OrderStatusFilter selected={selectedStatus} />

        {orders.length === 0 ? (
          <p className="text-gray-500 text-center py-12">No orders found.</p>
        ) : (
          <div className="space-y-3">
            {orders.map((order) => (
              <Link
                key={order.id}
                href={`/orders/${order.id}`}
                className="block bg-white rounded-lg border border-gray-200 p-4 shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <span className="font-bold">Order #{order.id}</span>
                    {order.customer && (
                      <span className="text-sm text-gray-500 ml-2">
                        {order.customer.name}
                      </span>
                    )}
                  </div>
                  <span
                    className={`text-xs font-medium px-2.5 py-1 rounded-full ${
                      statusColors[order.status] || "bg-gray-100 text-gray-700"
                    }`}
                  >
                    {order.status}
                  </span>
                </div>
                <div className="flex items-center justify-between mt-2 text-sm text-gray-500">
                  <span>${order.totalPrice.toFixed(2)}</span>
                  <span>
                    {new Date(order.createdAt).toLocaleDateString(undefined, {
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
