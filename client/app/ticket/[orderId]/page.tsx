import Link from "next/link";

type Ticket = {
  id: number;
  order_id: number;
  status: string;
  created_at: string;
  updated_at: string;
};

const KITCHEN_API_URL = process.env.KITCHEN_API_URL || "http://kitchen:8081";

const statusColors: Record<string, string> = {
  open: "bg-amber-100 text-amber-700",
  cooking: "bg-orange-100 text-orange-700",
  done: "bg-green-100 text-green-700",
};

export default async function TicketPage({
  params,
}: {
  params: Promise<{ orderId: string }>;
}) {
  const { orderId } = await params;

  let tickets: Ticket[] = [];
  let error: string | null = null;

  try {
    const res = await fetch(`${KITCHEN_API_URL}/ticket/${orderId}`, {
      cache: "no-store",
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    tickets = await res.json();
  } catch (e) {
    error = e instanceof Error ? e.message : "Failed to load tickets";
  }

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 py-8">
      <Link
        href={`/orders/${orderId}`}
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
        Order #{orderId}
      </Link>

      <h1 className="text-2xl font-bold mb-6">
        Tickets for Order #{orderId}
      </h1>

      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm mb-4">
          Could not load tickets. Make sure the kitchen service is running.
        </div>
      )}

      {!error && tickets.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500">No tickets found for this order.</p>
        </div>
      )}

      <div className="space-y-3">
        {tickets.map((ticket) => (
          <div
            key={ticket.id}
            className="bg-white rounded-lg border border-gray-200 p-4 shadow-sm"
          >
            <div className="flex items-center justify-between mb-2">
              <span className="font-bold">Ticket #{ticket.id}</span>
              <span
                className={`text-xs font-medium px-2.5 py-1 rounded-full ${
                  statusColors[ticket.status] || "bg-gray-100 text-gray-700"
                }`}
              >
                {ticket.status}
              </span>
            </div>
            <div className="flex items-center justify-between text-xs text-gray-400">
              <span>
                Created{" "}
                {new Date(ticket.created_at).toLocaleString(undefined, {
                  dateStyle: "medium",
                  timeStyle: "short",
                })}
              </span>
              {ticket.updated_at !== ticket.created_at && (
                <span>
                  Updated{" "}
                  {new Date(ticket.updated_at).toLocaleString(undefined, {
                    dateStyle: "medium",
                    timeStyle: "short",
                  })}
                </span>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
