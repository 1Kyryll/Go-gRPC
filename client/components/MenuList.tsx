"use client";

import { MenuItem } from "@/lib/types";
import { useOrder } from "./OrderContext";

export default function MenuList({ items }: { items: MenuItem[] }) {
  const { addItem, decrementItem, items: cartItems } = useOrder();

  function getCartQuantity(itemId: string) {
    return cartItems.find((i) => i.id === itemId)?.quantity ?? 0;
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {items.map((item) => {
        const qty = getCartQuantity(item.id);
        return (
          <div
            key={item.id}
            className={`rounded-lg border bg-white p-4 shadow-sm transition-all ${
              !item.isAvailable
                ? "opacity-50 border-gray-200"
                : qty > 0
                ? "border-indigo-300 ring-1 ring-indigo-200"
                : "border-gray-200 hover:shadow-md"
            }`}
          >
            <div className="flex items-start justify-between mb-1">
              <h2 className="text-base font-bold leading-tight">
                {item.name}
              </h2>
              <span className="text-base font-bold text-gray-900 shrink-0 ml-2">
                ${item.price.toFixed(2)}
              </span>
            </div>

            <p className="text-sm text-gray-500 mb-2 line-clamp-2">
              {item.description}
            </p>

            <div className="flex flex-wrap gap-1.5 mb-3">
              {!item.isAvailable && (
                <span className="text-xs font-medium px-2 py-0.5 rounded-full bg-red-100 text-red-700">
                  Sold out
                </span>
              )}
              {item.containsAllergens &&
                item.containsAllergens.length > 0 && (
                  <span
                    className="text-xs font-medium px-2 py-0.5 rounded-full bg-orange-100 text-orange-700"
                    title={item.containsAllergens.join(", ")}
                  >
                    Allergens
                  </span>
                )}
              {item.isAlcoholic && (
                <span className="text-xs font-medium px-2 py-0.5 rounded-full bg-purple-100 text-purple-700">
                  Alcoholic
                </span>
              )}
            </div>

            {item.isAvailable && (
              <div className="flex items-center justify-between">
                {qty === 0 ? (
                  <button
                    onClick={() =>
                      addItem({
                        id: item.id,
                        menuItemId: item.id,
                        menuItemName: item.name,
                        price: item.price,
                      })
                    }
                    className="w-full py-1.5 text-sm font-medium rounded-md bg-indigo-600 text-white hover:bg-indigo-700 transition-colors"
                  >
                    Add to Cart
                  </button>
                ) : (
                  <div className="flex items-center justify-between w-full">
                    <button
                      onClick={() => decrementItem(item.id)}
                      className="w-8 h-8 rounded-full border border-gray-300 flex items-center justify-center text-gray-600 hover:bg-gray-100 transition-colors text-lg leading-none"
                    >
                      -
                    </button>
                    <span className="text-sm font-bold">{qty} in cart</span>
                    <button
                      onClick={() =>
                        addItem({
                          id: item.id,
                          menuItemId: item.id,
                          menuItemName: item.name,
                          price: item.price,
                        })
                      }
                      className="w-8 h-8 rounded-full border border-gray-300 flex items-center justify-center text-gray-600 hover:bg-gray-100 transition-colors text-lg leading-none"
                    >
                      +
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
