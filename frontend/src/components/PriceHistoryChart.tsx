import { useEffect, useState } from "react";
import {
  LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid,
} from "recharts";

interface HistoryEntry {
  Date: string;
  "Close": number;
}

export default function PriceHistoryChart() {
  const [data, setData] = useState<HistoryEntry[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://backend:3001/api/history")
      .then((res) => res.json())
      .then((json) => {
        // Convert to proper format and reverse if necessary (to make oldest first)
        const cleaned = json.map((entry: any) => ({
          Date: entry.Date,
          Close: Number(entry["Close"]),
        })).reverse();
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
      <h2 className="text-2xl font-semibold text-gray-800 mb-4">Closing Price History</h2>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="Date" />
          <YAxis dataKey="Close" domain={["auto", "auto"]} />
          <Tooltip />
          <Line type="monotone" dataKey="Close" stroke="#3b82f6" strokeWidth={2} dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}