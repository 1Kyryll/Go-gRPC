import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

export async function proxy(request: NextRequest) {
    const token = request.cookies.get("token")?.value;

    const isProtectedRoute = request.nextUrl.pathname.startsWith("/orders") 
        || request.nextUrl.pathname.startsWith("/kitchen")
        || request.nextUrl.pathname.startsWith("/ticket")
        || request.nextUrl.pathname.startsWith("/order");

    if (isProtectedRoute && !token) {
        return NextResponse.redirect(new URL("/login", request.url));
    }

    return NextResponse.next();
}

export async function getCurrentUser() {
    const cookiesStore = await cookies();
    const user = cookiesStore.get("user")?.value; 

    if (!user) return;

    const decoded = decodeURIComponent(user);
    const parsed = JSON.parse(decoded); 

    return parsed; 
}

export const config = {
    matcher: ["/orders/:path*", "/kitchen/:path*", "/ticket/:path*", "/order/:path*"],
};