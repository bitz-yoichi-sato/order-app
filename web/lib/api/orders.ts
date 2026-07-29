export type Order = {
  id: string;
  customer: string;
  status: string;
  amount: number;
  orderedAt: string;
};

export type OrdersResponse = {
  orders: Order[];
  total: number;
};

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

export async function fetchOrders(): Promise<OrdersResponse> {
  const res = await fetch(`${API_BASE}/api/orders`, { cache: "no-store" });
  if (!res.ok) {
    throw new Error(`failed to fetch orders: ${res.status}`);
  }
  return res.json();
}
