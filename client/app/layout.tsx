import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import "./globals.css";
import { OrderProvider } from "../components/OrderContext";

const roboto = Roboto({
  subsets: ["latin"],
  variable: "--font-roboto",
  weight: ["400", "500", "700"],
});

export const metadata: Metadata = {
  title: "OrderFlow",
  description: "Order & kitchen management system",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${roboto.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">
        <OrderProvider>
          <main className="flex-1">
            {children}
          </main>
        </OrderProvider>
      </body>
    </html>
  );
}
