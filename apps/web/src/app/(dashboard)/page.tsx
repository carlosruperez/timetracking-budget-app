"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useTimer } from "@/hooks/useTimer";
import { ActiveTimer } from "@/components/timer/ActiveTimer";
import { BudgetStatusList } from "@/components/budget/BudgetStatus";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardContent } from "@/components/ui/card";

export default function DashboardPage() {
  const { start } = useTimer();
  const [selectedCategory, setSelectedCategory] = useState("");
  const [description, setDescription] = useState("");

  const { data: categories } = useQuery({
    queryKey: ["categories"],
    queryFn: () => api.listCategories(),
  });

  const handleStart = () => {
    if (!selectedCategory) return;
    start.mutate({ categoryId: selectedCategory, description });
    setDescription("");
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>

      <ActiveTimer />

      <Card>
        <CardHeader>
          <h2 className="font-semibold text-gray-900">Start Timer</h2>
        </CardHeader>
        <CardContent>
          <div className="flex gap-3">
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
            >
              <option value="">Select category...</option>
              {categories?.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
            <input
              type="text"
              placeholder="Description (optional)"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none"
            />
            <Button
              onClick={handleStart}
              disabled={!selectedCategory || start.isPending}
            >
              Start
            </Button>
          </div>
        </CardContent>
      </Card>

      <div>
        <h2 className="font-semibold text-gray-900 mb-3">Budget Status</h2>
        <BudgetStatusList />
      </div>
    </div>
  );
}
