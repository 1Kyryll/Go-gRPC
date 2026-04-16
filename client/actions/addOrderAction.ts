"use server";

import { OrderItem } from "@/lib/types";
import { gql } from "@/lib/graphql";
import { getCurrentUser } from "@/proxy";

type CreateOrderResult = {
  success: boolean;
  orderId?: number;
  errors?: { field: string; message: string }[];
};

export const addOrderAction = async (
  formData: FormData,
  items: OrderItem[]
): Promise<CreateOrderResult> => {
  const user = await getCurrentUser(); 

  const user_id = user?.id;

  if (!user_id) {
    return { success: false, errors: [{ field: "userId", message: "Unauthorized user" }] };
  }

  if (items.length === 0) {
    return { success: false, errors: [{ field: "items", message: "Cart is empty" }] };
  }

  try {
    const data = await gql<{
      createOrder: {
        order: { id: number } | null;
        errors: { field: string; message: string }[] | null;
      };
    }>(
      `mutation CreateOrder($input: CreateOrderInput!) {
        createOrder(input: $input) {
            order { id }
            errors { field, message }
        }
      }`,
      {
        input: {
          userId: parseInt(user_id),
          items: items.map((item) => ({
            menuItemId: item.menuItemId,
            quantity: item.quantity,
            specialInstructions: item.specialInstructions || "",
          })),
        },
      }
    );

    if (data.createOrder.errors && data.createOrder.errors.length > 0) {
      return { success: false, errors: data.createOrder.errors };
    }

    return { success: true, orderId: data.createOrder.order?.id };
  } catch (err) {
    return {
      success: false,
      errors: [{ field: "general", message: err instanceof Error ? err.message : "Failed to create order" }],
    };
  }
};
