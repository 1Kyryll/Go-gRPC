export type MenuItem = {
    id: string;
    name: string;
    description: string;
    price: number;
    isAvailable: boolean;
    category: string;
    containsAllergens?: string[];
    isAlcoholic?: boolean;
};

export type OrderItem = {
    id: string;
    menuItemId: string;
    menuItemName: string;
    quantity: number;
    specialInstructions?: string;
    price: number;
}
