import { gql } from "../lib/graphql";
import { MenuItem } from "../lib/types";
import MenuList from "../components/MenuList";
import CategoryFilter from "../components/CategoryFilter";
import Navbar from "@/components/Navbar";

type Props = {
  searchParams: Promise<{ category?: string }>;
};

export default async function Home({ searchParams }: Props) {
  const { category } = await searchParams;

  const selectedCategory =
    category && ["APPETIZER", "MAIN", "DRINK", "DESSERT"].includes(category)
      ? category
      : null;

  const data = await gql<{
    menuItems: { edges: { node: MenuItem }[]; totalCount: number };
  }>(
    `query GetMenu($category: MenuCategory, $first: Int) {
      menuItems(category: $category, first: $first) {
        edges {
          node {
            ... on FoodItem {
              id, name, description, price, isAvailable, containsAllergens
            }
            ... on DrinkItem {
              id, name, description, price, isAvailable, isAlcoholic
            }
          }
        }
        totalCount
      }
    }`,
    { category: selectedCategory, first: 100 }
  );

  const items = data.menuItems.edges.map((edge) => edge.node);

  return (
    <>
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6">
          <h1 className="text-2xl font-bold mb-1">Menu</h1>
          <p className="text-sm text-gray-500">
            {data.menuItems.totalCount} items available
          </p>
        </div>

        <CategoryFilter selected={selectedCategory} />

        <MenuList items={items} />
      </div>
    </>
  );
}
