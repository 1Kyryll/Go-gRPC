export default function Home() {
  return (
    <div className="p-8">
      <h1 className="text-4xl font-bold mb-6">Welcome</h1>
      <a href="/order/create" className="text-blue-500 hover:underline">Create Order</a>
      <a href="/kitchen" className="ml-4 text-blue-500 hover:underline">Kitchen View</a>
    </div>
  );
}
