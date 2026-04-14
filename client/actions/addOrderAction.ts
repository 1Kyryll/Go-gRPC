"use server";

import { OrderItem } from "@/lib/types";
import { gql } from "@/lib/graphql";

type CreateOrderResult = {
  success: boolean;
  orderId?: number;
  errors?: { field: string; message: string }[];
};

export const addOrderAction = async (
  formData: FormData,
  items: OrderItem[]
): Promise<CreateOrderResult> => {
  const customer_id = formData.get("customer_id") as string;

  if (!customer_id) {
    return { success: false, errors: [{ field: "customerId", message: "Customer ID is required" }] };
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
          customerId: parseInt(customer_id),
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
