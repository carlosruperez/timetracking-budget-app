import { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  TextInput,
  Alert,
} from "react-native";
import { useQuery } from "@tanstack/react-query";
import { SafeAreaView } from "react-native-safe-area-context";
import { useTimer } from "@/hooks/useTimer";
import { api } from "@/lib/api";
import { formatDuration } from "@timetracking/core";

export default function TimerScreen() {
  const { activeTimer, elapsed, start, pause, resume, stop } = useTimer();
  const [selectedCategory, setSelectedCategory] = useState("");
  const [description, setDescription] = useState("");

  const { data: categories = [] } = useQuery({
    queryKey: ["categories"],
    queryFn: () => api.listCategories(),
  });

  const handleStart = () => {
    if (!selectedCategory) {
      Alert.alert("Select a category first");
      return;
    }
    start.mutate({ categoryId: selectedCategory, description });
    setDescription("");
  };

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <ScrollView contentContainerStyle={styles.scroll}>
        {/* Active timer display */}
        <View style={styles.timerCard}>
          {activeTimer ? (
            <>
              <Text style={styles.timerDisplay}>{formatDuration(elapsed)}</Text>
              <View style={styles.timerMeta}>
                <View
                  style={[
                    styles.dot,
                    {
                      backgroundColor:
                        activeTimer.state === "running" ? "#22c55e" : "#eab308",
                    },
                  ]}
                />
                <Text style={styles.timerState}>{activeTimer.state}</Text>
              </View>
              {activeTimer.description ? (
                <Text style={styles.timerDesc}>{activeTimer.description}</Text>
              ) : null}

              {/* Controls */}
              <View style={styles.controls}>
                {activeTimer.state === "running" ? (
                  <TouchableOpacity
                    style={[styles.controlBtn, styles.pauseBtn]}
                    onPress={() => pause.mutate()}
                    disabled={pause.isPending}
                  >
                    <Text style={[styles.controlBtnText, styles.pauseBtnText]}>⏸ Pause</Text>
                  </TouchableOpacity>
                ) : (
                  <TouchableOpacity
                    style={[styles.controlBtn, styles.resumeBtn]}
                    onPress={() => resume.mutate()}
                    disabled={resume.isPending}
                  >
                    <Text style={styles.controlBtnText}>▶ Resume</Text>
                  </TouchableOpacity>
                )}
                <TouchableOpacity
                  style={[styles.controlBtn, styles.stopBtn]}
                  onPress={() => stop.mutate()}
                  disabled={stop.isPending}
                >
                  <Text style={styles.controlBtnText}>⏹ Stop</Text>
                </TouchableOpacity>
              </View>
            </>
          ) : (
            <Text style={styles.noTimer}>No active timer</Text>
          )}
        </View>

        {/* Start timer */}
        {!activeTimer && (
          <View style={styles.startCard}>
            <Text style={styles.sectionTitle}>Start Timer</Text>

            <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.categoryScroll}>
              {categories.map((cat) => (
                <TouchableOpacity
                  key={cat.id}
                  style={[
                    styles.categoryChip,
                    selectedCategory === cat.id && styles.categoryChipSelected,
                    { borderColor: cat.color },
                  ]}
                  onPress={() => setSelectedCategory(cat.id)}
                >
                  <View style={[styles.categoryDot, { backgroundColor: cat.color }]} />
                  <Text
                    style={[
                      styles.categoryChipText,
                      selectedCategory === cat.id && { color: cat.color },
                    ]}
                  >
                    {cat.name}
                  </Text>
                </TouchableOpacity>
              ))}
            </ScrollView>

            <TextInput
              style={styles.input}
              placeholder="Description (optional)"
              value={description}
              onChangeText={setDescription}
            />

            <TouchableOpacity
              style={[
                styles.startBtn,
                !selectedCategory && styles.startBtnDisabled,
              ]}
              onPress={handleStart}
              disabled={!selectedCategory || start.isPending}
            >
              <Text style={styles.startBtnText}>
                {start.isPending ? "Starting..." : "▶ Start Timer"}
              </Text>
            </TouchableOpacity>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#f8fafc" },
  scroll: { padding: 16, gap: 16 },
  timerCard: {
    backgroundColor: "#fff",
    borderRadius: 16,
    padding: 24,
    alignItems: "center",
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 8,
    elevation: 2,
  },
  timerDisplay: {
    fontSize: 56,
    fontWeight: "bold",
    fontVariant: ["tabular-nums"],
    color: "#1e293b",
    letterSpacing: 2,
  },
  timerMeta: { flexDirection: "row", alignItems: "center", gap: 6, marginTop: 8 },
  dot: { width: 8, height: 8, borderRadius: 4 },
  timerState: { fontSize: 14, color: "#64748b", textTransform: "capitalize" },
  timerDesc: { fontSize: 14, color: "#94a3b8", marginTop: 4 },
  controls: { flexDirection: "row", gap: 12, marginTop: 20 },
  controlBtn: {
    paddingHorizontal: 20,
    paddingVertical: 12,
    borderRadius: 10,
    minWidth: 110,
    alignItems: "center",
  },
  pauseBtn: { backgroundColor: "#f1f5f9", borderWidth: 1, borderColor: "#e2e8f0" },
  resumeBtn: { backgroundColor: "#6366f1" },
  stopBtn: { backgroundColor: "#ef4444" },
  controlBtnText: { fontWeight: "600", fontSize: 15, color: "#fff" },
  pauseBtnText: { color: "#475569" },
  noTimer: { fontSize: 16, color: "#94a3b8", paddingVertical: 16 },
  startCard: {
    backgroundColor: "#fff",
    borderRadius: 16,
    padding: 20,
    gap: 12,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 8,
    elevation: 2,
  },
  sectionTitle: { fontSize: 17, fontWeight: "600", color: "#1e293b" },
  categoryScroll: { marginHorizontal: -4 },
  categoryChip: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 20,
    borderWidth: 1.5,
    borderColor: "#e2e8f0",
    marginHorizontal: 4,
    backgroundColor: "#f8fafc",
  },
  categoryChipSelected: { backgroundColor: "#f0f9ff" },
  categoryDot: { width: 8, height: 8, borderRadius: 4 },
  categoryChipText: { fontSize: 14, color: "#475569", fontWeight: "500" },
  input: {
    backgroundColor: "#f8fafc",
    borderWidth: 1,
    borderColor: "#e2e8f0",
    borderRadius: 10,
    padding: 12,
    fontSize: 15,
  },
  startBtn: {
    backgroundColor: "#6366f1",
    borderRadius: 10,
    padding: 16,
    alignItems: "center",
  },
  startBtnDisabled: { opacity: 0.4 },
  startBtnText: { color: "#fff", fontSize: 16, fontWeight: "600" },
});
