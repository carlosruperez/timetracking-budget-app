import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  FlatList,
  Alert,
  RefreshControl,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

const PRESET_COLORS = [
  "#6366f1", "#ec4899", "#f97316", "#22c55e",
  "#06b6d4", "#8b5cf6", "#eab308", "#ef4444",
];

export default function CategoriesScreen() {
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [color, setColor] = useState("#6366f1");

  const { data: categories = [], isLoading, refetch } = useQuery({
    queryKey: ["categories"],
    queryFn: () => api.listCategories(),
  });

  const create = useMutation({
    mutationFn: () => api.createCategory(name, color),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["categories"] });
      setName("");
    },
    onError: (err: Error) => Alert.alert("Error", err.message),
  });

  const remove = useMutation({
    mutationFn: (id: string) => api.deleteCategory(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["categories"] }),
    onError: (err: Error) => Alert.alert("Error", err.message),
  });

  const handleDelete = (id: string, catName: string) => {
    Alert.alert("Delete Category", `Delete "${catName}"?`, [
      { text: "Cancel", style: "cancel" },
      { text: "Delete", style: "destructive", onPress: () => remove.mutate(id) },
    ]);
  };

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      {/* Create form */}
      <View style={styles.form}>
        <TextInput
          style={styles.input}
          placeholder="Category name"
          value={name}
          onChangeText={setName}
        />
        <View style={styles.colorRow}>
          {PRESET_COLORS.map((c) => (
            <TouchableOpacity
              key={c}
              style={[
                styles.colorSwatch,
                { backgroundColor: c },
                color === c && styles.colorSwatchSelected,
              ]}
              onPress={() => setColor(c)}
            />
          ))}
        </View>
        <TouchableOpacity
          style={[styles.createBtn, !name && styles.createBtnDisabled]}
          onPress={() => create.mutate()}
          disabled={!name || create.isPending}
        >
          <Text style={styles.createBtnText}>
            {create.isPending ? "Adding..." : "+ Add Category"}
          </Text>
        </TouchableOpacity>
      </View>

      <FlatList
        data={categories}
        keyExtractor={(item) => item.id}
        contentContainerStyle={styles.list}
        refreshControl={
          <RefreshControl refreshing={isLoading} onRefresh={refetch} />
        }
        renderItem={({ item }) => (
          <View style={styles.categoryRow}>
            <View style={[styles.dot, { backgroundColor: item.color }]} />
            <Text style={styles.categoryName}>{item.name}</Text>
            <TouchableOpacity
              onPress={() => handleDelete(item.id, item.name)}
              hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
            >
              <Text style={styles.deleteText}>Delete</Text>
            </TouchableOpacity>
          </View>
        )}
        ListEmptyComponent={
          !isLoading ? (
            <Text style={styles.emptyText}>No categories yet.</Text>
          ) : null
        }
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: "#f8fafc" },
  form: {
    backgroundColor: "#fff",
    padding: 16,
    gap: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#e2e8f0",
  },
  input: {
    backgroundColor: "#f8fafc",
    borderWidth: 1,
    borderColor: "#e2e8f0",
    borderRadius: 10,
    padding: 12,
    fontSize: 15,
  },
  colorRow: { flexDirection: "row", gap: 10, flexWrap: "wrap" },
  colorSwatch: {
    width: 28,
    height: 28,
    borderRadius: 14,
  },
  colorSwatchSelected: {
    borderWidth: 3,
    borderColor: "#1e293b",
  },
  createBtn: {
    backgroundColor: "#6366f1",
    borderRadius: 10,
    padding: 14,
    alignItems: "center",
  },
  createBtnDisabled: { opacity: 0.4 },
  createBtnText: { color: "#fff", fontWeight: "600", fontSize: 15 },
  list: { padding: 16, gap: 8 },
  categoryRow: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: "#fff",
    borderRadius: 12,
    padding: 14,
    gap: 12,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.04,
    shadowRadius: 4,
    elevation: 1,
  },
  dot: { width: 14, height: 14, borderRadius: 7 },
  categoryName: { flex: 1, fontSize: 15, color: "#1e293b", fontWeight: "500" },
  deleteText: { fontSize: 14, color: "#ef4444" },
  emptyText: { textAlign: "center", color: "#94a3b8", marginTop: 32, fontSize: 15 },
});
