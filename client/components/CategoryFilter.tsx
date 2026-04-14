"use client";

import { useRouter, useSearchParams } from "next/navigation";

const categories = [
  { value: null, label: "All" },
  { value: "APPETIZER", label: "Appetizers" },
  { value: "MAIN", label: "Mains" },
  { value: "DRINK", label: "Drinks" },
  { value: "DESSERT", label: "Desserts" },
];

export default function CategoryFilter({
  selected,
}: {
  selected: string | null;
}) {
  const router = useRouter();
  const searchParams = useSearchParams();

  function handleSelect(value: string | null) {
    const params = new URLSearchParams(searchParams.toString());
    if (value) {
      params.set("category", value);
    } else {
      params.delete("category");
    }
    router.push(`/?${params.toString()}`);
  }

  return (
    <div className="flex gap-2 mb-6 overflow-x-auto pb-1">
      {categories.map((cat) => {
        const active = selected === cat.value;
        return (
          <button
            key={cat.label}
            onClick={() => handleSelect(cat.value)}
            className={`px-3 py-1.5 text-sm font-medium rounded-full whitespace-nowrap transition-colors ${
              active
                ? "bg-emerald-600 text-white"
                : "bg-white border border-gray-200 text-gray-600 hover:bg-gray-50"
            }`}
          >
            {cat.label}
          </button>
        );
      })}
    </div>
  );
}
