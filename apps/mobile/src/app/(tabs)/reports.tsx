import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { formatDuration } from "@timetracking/core";

function getLastWeekRange() {
  const to = new Date();
  const from = new Date();
  from.setDate(from.getDate() - 7);
  return { from: from.toISOString(), to: to.toISOString() };
}

export default function ReportsScreen() {
  const { from, to } = getLastWeekRange();

  const { data: summary = [], isLoading, refetch } = useQuery({
    queryKey: ["reports", "summary", from, to],
    queryFn: () => api.getSummary(from, to),
  });

  const totalSec = summary.reduce((acc, s) => acc + s.total_sec, 0);
  const maxSec = summary.length > 0 ? Math.max(...summary.map((s) => s.total_sec)) : 1;

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <ScrollView
        contentContainerStyle={styles.scroll}
        refreshControl={
          <RefreshControl refreshing={isLoading} onRefresh={refetch} />
        }
      >
        <View style={styles.totalCard}>
          <Text style={styles.totalLabel}>Total this week</Text>
          <Text style={styles.totalValue}>{formatDuration(totalSec)}</Text>
        </View>

        <Text style={styles.sectionTitle}>By Category</Text>

        {summary.length === 0 && !isLoading ? (
          <Text style={styles.emptyText}>No tracked time this week.</Text>
        ) : (
          summary.map((item) => {
            const pct = totalSec > 0 ? (item.total_sec / maxSec) * 100 : 0;
            return (
              <View key={item.category_id} style={styles.row}>
                <View style={styles.rowHeader}>
                  <Text style={styles.categoryName}>{item.category_name}</Text>
                  <Text style={styles.duration}>{formatDuration(item.total_sec)}</Text>
                </View>
                <View style={styles.barTrack}>
                  <View
                    style={[styles.bar, { width: `${pct}%`, backgroundColor: "#6366f1" }]}
                  />
                </View>
                <Text style={styles.sessions}>{item.entry_count} sessions</Text>
              </View>
            );
          })
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#f8fafc" },
  scroll: { padding: 16, gap: 16 },
  totalCard: {
    backgroundColor: "#6366f1",
    borderRadius: 16,
    padding: 24,
    alignItems: "center",
  },
  totalLabel: { color: "#c7d2fe", fontSize: 14, fontWeight: "500" },
  totalValue: {
    color: "#fff",
    fontSize: 40,
    fontWeight: "bold",
    fontVariant: ["tabular-nums"],
    marginTop: 4,
  },
  sectionTitle: { fontSize: 17, fontWeight: "600", color: "#1e293b" },
  emptyText: { textAlign: "center", color: "#94a3b8", fontSize: 15, marginTop: 16 },
  row: {
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 14,
    gap: 6,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.04,
    shadowRadius: 4,
    elevation: 1,
  },
  rowHeader: { flexDirection: "row", justifyContent: "space-between", alignItems: "center" },
  categoryName: { fontSize: 15, fontWeight: "600", color: "#1e293b" },
  duration: { fontSize: 14, fontWeight: "600", fontVariant: ["tabular-nums"], color: "#6366f1" },
  barTrack: { height: 6, backgroundColor: "#f1f5f9", borderRadius: 3, overflow: "hidden" },
  bar: { height: 6, borderRadius: 3 },
  sessions: { fontSize: 12, color: "#94a3b8" },
});
