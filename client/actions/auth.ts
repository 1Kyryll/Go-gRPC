"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { gql } from "@/lib/graphql";

type AuthResult = {
  success: boolean;
  error?: string;
};

type AuthResponse = {
  authToken: {
    accessToken: string;
    expiredAt: string;
  };
  user: {
    id: string;
    username: string;
    email: string;
    phone?: string;
  };
};

export async function loginAction(formData: FormData): Promise<AuthResult> {
  const username = formData.get("username") as string;
  const password = formData.get("password") as string;

  if (!username || !password) {
    return { success: false, error: "Username and password are required" };
  }

  try {
    const data = await gql<{ login: AuthResponse }>(
      `mutation Login($input: LoginInput!) {
        login(input: $input) {
          authToken {
            accessToken
            expiredAt
          }
          user {
            id
            username
            email
            phone
          }
        }
      }`,
      { input: { username, password } }
    );

    const cookieStore = await cookies();
    cookieStore.set("token", data.login.authToken.accessToken, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      expires: new Date(data.login.authToken.expiredAt),
    });
    cookieStore.set("user", JSON.stringify(data.login.user), {
      httpOnly: false,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      expires: new Date(data.login.authToken.expiredAt),
    });
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err.message : "Login failed",
    };
  }

  redirect("/");
}

export async function registerAction(formData: FormData): Promise<AuthResult> {
  const username = formData.get("username") as string;
  const email = formData.get("email") as string;
  const password = formData.get("password") as string;
  const phone = (formData.get("phone") as string) || undefined;

  if (!username || !email || !password) {
    return { success: false, error: "Username, email, and password are required" };
  }

  try {
    const data = await gql<{ register: AuthResponse }>(
      `mutation Register($input: RegisterInput!) {
        register(input: $input) {
          authToken {
            accessToken
            expiredAt
          }
          user {
            id
            username
            email
            phone
          }
        }
      }`,
      { input: { username, email, password, phone } }
    );

    const cookieStore = await cookies();
    cookieStore.set("token", data.register.authToken.accessToken, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      expires: new Date(data.register.authToken.expiredAt),
    });
    cookieStore.set("user", JSON.stringify(data.register.user), {
      httpOnly: false,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      expires: new Date(data.register.authToken.expiredAt),
    });
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err.message : "Registration failed",
    };
  }

  redirect("/");
}

export async function logoutAction(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete("token");
  cookieStore.delete("user");
  redirect("/login");
}
