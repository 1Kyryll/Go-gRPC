type Ticket = {
    id: number;
    order_id: number;
    status: string;
    created_at: string;
    updated_at: string;
};

export default async function TicketPage({ params }: { params: Promise<{ orderId: string }> }) {
    const { orderId } = await params;
    const res = await fetch(`http://localhost:8080/ticket/${orderId}`, { cache: "no-store" });
    const tickets: Ticket[] = await res.json();

    return (
        <div className="p-8">
            <h1 className="text-3xl font-bold mb-6">Tickets for Order #{orderId}</h1>

            {tickets.length === 0 && (
                <p className="text-gray-500">No tickets found for this order.</p>
            )}

            <div className="space-y-4">
                {tickets.map((ticket) => (
                    <div
                        key={ticket.id}
                        className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
                    >
                        <div className="flex justify-between items-center mb-2">
                            <span className="font-bold text-lg">Ticket #{ticket.id}</span>
                            <span className={`text-sm font-medium px-2 py-1 rounded ${
                                ticket.status === "open"
                                    ? "bg-yellow-100 text-yellow-800"
                                    : "bg-green-100 text-green-800"
                            }`}>
                                {ticket.status}
                            </span>
                        </div>
                        <p className="text-sm text-gray-500">Order #{ticket.order_id}</p>
                    </div>
                ))}
            </div>
        </div>
    );
}
