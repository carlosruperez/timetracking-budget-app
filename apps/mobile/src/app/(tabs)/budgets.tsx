import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  RefreshControl,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useBudgets } from "@/hooks/useBudgets";
import { getBudgetAlert, formatBudgetSeconds } from "@timetracking/core";

export default function BudgetsScreen() {
  const { data: statuses = [], isLoading, refetch } = useBudgets();

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <ScrollView
        contentContainerStyle={styles.scroll}
        refreshControl={
          <RefreshControl refreshing={isLoading} onRefresh={refetch} />
        }
      >
        {statuses.length === 0 && !isLoading ? (
          <View style={styles.empty}>
            <Text style={styles.emptyText}>
              No budget rules set.{"\n"}Create rules via the web app.
            </Text>
          </View>
        ) : (
          statuses.map((status) => {
            const alert = getBudgetAlert(status);
            const pct = Math.min(100, status.percent_used);
            const barColor =
              alert === "exceeded"
                ? "#ef4444"
                : alert === "warning"
                ? "#eab308"
                : "#22c55e";
            const badgeBg =
              alert === "exceeded"
                ? "#fee2e2"
                : alert === "warning"
                ? "#fef9c3"
                : "#dcfce7";
            const badgeText =
              alert === "exceeded"
                ? "#991b1b"
                : alert === "warning"
                ? "#854d0e"
                : "#166534";

            return (
              <View key={status.rule.id} style={styles.card}>
                <View style={styles.cardHeader}>
                  <Text style={styles.periodLabel}>
                    {status.rule.period_type.charAt(0).toUpperCase() +
                      status.rule.period_type.slice(1)}{" "}
                    Budget
                  </Text>
                  <View style={[styles.badge, { backgroundColor: badgeBg }]}>
                    <Text style={[styles.badgeText, { color: badgeText }]}>
                      {Math.round(status.percent_used)}%
                    </Text>
                  </View>
                </View>

                <View style={styles.progressTrack}>
                  <View
                    style={[
                      styles.progressBar,
                      { width: `${pct}%`, backgroundColor: barColor },
                    ]}
                  />
                </View>

                <View style={styles.cardFooter}>
                  <Text style={styles.footerText}>
                    {formatBudgetSeconds(status.used_sec)} used
                  </Text>
                  <Text style={styles.footerText}>
                    {formatBudgetSeconds(status.rule.budget_sec)} budget
                  </Text>
                </View>

                <Text style={styles.remainingText}>
                  {alert === "exceeded"
                    ? `Over by ${formatBudgetSeconds(status.used_sec - status.rule.budget_sec)}`
                    : `${formatBudgetSeconds(status.remaining_sec)} remaining`}
                </Text>
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
  scroll: { padding: 16, gap: 12 },
  empty: { alignItems: "center", marginTop: 64 },
  emptyText: { fontSize: 15, color: "#94a3b8", textAlign: "center", lineHeight: 22 },
  card: {
    backgroundColor: "#fff",
    borderRadius: 16,
    padding: 18,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 8,
    elevation: 2,
    gap: 10,
  },
  cardHeader: { flexDirection: "row", justifyContent: "space-between", alignItems: "center" },
  periodLabel: { fontSize: 16, fontWeight: "600", color: "#1e293b" },
  badge: { paddingHorizontal: 10, paddingVertical: 4, borderRadius: 12 },
  badgeText: { fontSize: 12, fontWeight: "600" },
  progressTrack: {
    height: 8,
    backgroundColor: "#f1f5f9",
    borderRadius: 4,
    overflow: "hidden",
  },
  progressBar: { height: 8, borderRadius: 4 },
  cardFooter: { flexDirection: "row", justifyContent: "space-between" },
  footerText: { fontSize: 12, color: "#94a3b8" },
  remainingText: { fontSize: 13, color: "#475569", fontWeight: "500" },
});
