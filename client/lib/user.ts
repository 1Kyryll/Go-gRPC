export function getClientUser(): { id: string; username: string; role: string; } | null {
    if (typeof document === "undefined") return null;
    const match = document.cookie.match(/user=([^;]+)/);
    if (!match) return null;
    try {
        return JSON.parse(decodeURIComponent(match[1]));
    } catch {
        return null;
    }
}

export function getClientToken(): string | null {
    if (typeof document === "undefined") return null;
    const match = document.cookie.match(/token=([^;]+)/);
    if (!match) return null;
    return decodeURIComponent(match[1]);
}