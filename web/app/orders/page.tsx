import { fetchOrders } from "@/lib/api/orders";

export default async function OrdersPage() {
  const { orders, total } = await fetchOrders();

  return (
    <main style={{ padding: "2rem" }}>
      <h1>注文一覧</h1>
      <p>全 {total} 件</p>
      <table border={1} cellPadding={8} style={{ borderCollapse: "collapse" }}>
        <thead>
          <tr>
            <th>ID</th>
            <th>顧客</th>
            <th>ステータス</th>
            <th>金額</th>
            <th>注文日時</th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>{order.id}</td>
              <td>{order.customer}</td>
              <td>{order.status}</td>
              <td>{order.amount}</td>
              <td>{order.orderedAt}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </main>
  );
}
