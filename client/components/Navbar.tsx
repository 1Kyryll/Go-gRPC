"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useOrder } from "./OrderContext";
import { getClientUser } from "@/lib/user";

export default function Navbar() {
  const { totalItems } = useOrder();
  const pathname = usePathname();
  const user = getClientUser(); 

  const links = [
    { href: "/", label: "Menu" },
    ...(user?.role != "KITCHEN_STAFF" ? [{ href: "/orders", label: "Orders" }] : []),
    ...(user?.role === "KITCHEN_STAFF" ? [{ href: "/kitchen", label: "Kitchen" }] : []),
  ];

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-14">
          <div className="flex items-center gap-1">
            <Link
              href="/"
              className="text-lg font-bold text-emerald-700 mr-6 shrink-0"
            >
              OrderFlow
            </Link>
            <div className="flex items-center gap-1">
              {links.map((link) => {
                const active = pathname === link.href;
                return (
                  <Link
                    key={link.href}
                    href={link.href}
                    className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
                      active
                        ? "bg-emerald-50 text-emerald-700"
                        : "text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                    }`}
                  >
                    {link.label}
                  </Link>
                );
              })}
            </div>
          </div>

          <div className="flex items-center gap-2">
            {!user && (
              <>
                <Link
                  href="/login"
                  className="px-3 py-1.5 rounded-md text-sm font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50 transition-colors"
                >
                  Login
                </Link>
                <Link
                  href="/signup"
                  className="px-3 py-1.5 rounded-md text-sm font-medium bg-emerald-600 text-white hover:bg-emerald-700 transition-colors"
                >
                  Sign up
                </Link>
              </>
            )}
            {user && user.role !== "KITCHEN_STAFF" && (
              <Link
                href="/order/create"
                className="relative inline-flex items-center gap-2 px-4 py-1.5 rounded-md text-sm font-medium bg-emerald-600 text-white hover:bg-emerald-700 transition-colors"
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
                    d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"
                  />
                </svg>
                Cart
                {totalItems > 0 && (
                  <span className="absolute -top-1.5 -right-1.5 inline-flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-red-500 rounded-full">
                    {totalItems > 99 ? "99+" : totalItems}
                  </span>
                )}
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
