"use client";

import { useRouter, useSearchParams } from "next/navigation";

const statuses = [
  { value: null, label: "All" },
  { value: "PENDING", label: "Pending" },
  { value: "CONFIRMED", label: "Confirmed" },
  { value: "COOKING", label: "Cooking" },
  { value: "COMPLETED", label: "Completed" },
  { value: "CANCELLED", label: "Cancelled" },
];

export default function OrderStatusFilter({
  selected,
}: {
  selected: string | null;
}) {
  const router = useRouter();
  const searchParams = useSearchParams();

  function handleSelect(value: string | null) {
    const params = new URLSearchParams(searchParams.toString());
    if (value) {
      params.set("status", value);
    } else {
      params.delete("status");
    }
    router.push(`/orders?${params.toString()}`);
  }

  return (
    <div className="flex gap-2 mb-6 overflow-x-auto pb-1">
      {statuses.map((s) => {
        const active = selected === s.value;
        return (
          <button
            key={s.label}
            onClick={() => handleSelect(s.value)}
            className={`px-3 py-1.5 text-sm font-medium rounded-full whitespace-nowrap transition-colors ${
              active
                ? "bg-indigo-600 text-white"
                : "bg-white border border-gray-200 text-gray-600 hover:bg-gray-50"
            }`}
          >
            {s.label}
          </button>
        );
      })}
    </div>
  );
}
