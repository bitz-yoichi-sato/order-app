import { fetchOrders, type Order } from "@/lib/api/orders";

const STATUS_LABEL: Record<string, string> = {
  pending: "保留中",
  shipped: "発送済み",
  canceled: "キャンセル",
};

const STATUS_BADGE_CLASS: Record<string, string> = {
  pending: "badge--pending",
  shipped: "badge--shipped",
  canceled: "badge--canceled",
};

const amountFormatter = new Intl.NumberFormat("ja-JP", {
  style: "currency",
  currency: "JPY",
});

const dateFormatter = new Intl.DateTimeFormat("ja-JP", {
  timeZone: "UTC",
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
});

function StatusBadge({ status }: { status: Order["status"] }) {
  const badgeClass = STATUS_BADGE_CLASS[status] ?? "badge--canceled";
  const label = STATUS_LABEL[status] ?? status;
  return <span className={`badge ${badgeClass}`}>{label}</span>;
}

export default async function OrdersPage() {
  const { orders, total } = await fetchOrders();

  return (
    <main className="page">
      <div className="page__eyebrow">Orders</div>
      <h1 className="page__title">注文一覧</h1>
      <p className="page__lead">全 {total} 件の注文を表示しています。</p>

      <section className="card">
        <div className="card__header">
          <span className="card__title">注文</span>
          <span className="card__meta">{orders.length} 件を表示</span>
        </div>
        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>注文ID</th>
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
                  <td>
                    <StatusBadge status={order.status} />
                  </td>
                  <td className="amount">
                    {amountFormatter.format(order.amount)}
                  </td>
                  <td>{dateFormatter.format(new Date(order.orderedAt))}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </main>
  );
}
