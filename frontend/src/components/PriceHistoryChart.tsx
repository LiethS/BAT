import { useEffect, useState } from "react";
import {
  LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid,
} from "recharts";

// Match the Go API response structure
interface HistoryEntry {
  date: string;
  price: number;
}

export default function PriceHistoryChart() {
  const [data, setData] = useState<HistoryEntry[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:3001/api/history")
      .then((res) => res.json())
      .then((json: HistoryEntry[]) => {
        // Sort by ascending date just in case (oldest → newest)
        const cleaned = [...json].sort(
          (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
        );
        setData(cleaned);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Failed to fetch price history:", err);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <p className="text-center text-gray-600">Loading chart...</p>;
  }

  return (
    <div className="w-full max-w-4xl mx-auto bg-white shadow-md rounded-xl p-6">
      <h2 className="text-2xl font-semibold text-gray-800 mb-4">Bitcoin Price History</h2>
      <ResponsiveContainer width="100%" height={700}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="date" />
          <YAxis dataKey="price" domain={["auto", "auto"]} />
          <Tooltip />
          <Line
            type="monotone"
            dataKey="price"
            stroke="#3b82f6"
            strokeWidth={2}
            dot={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}