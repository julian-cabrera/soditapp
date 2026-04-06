<script lang="ts">
  import { Input } from "$lib/components/ui/input/index.js";
  import {
    getWeekStart,
    getWeekEnd,
    setWeekRange,
  } from "$lib/stores/delivery.svelte.js";

  let weekStart = $derived(getWeekStart());
  let weekEnd = $derived(getWeekEnd());

  let inputValue = $derived(toInputValue(weekStart));

  let weekLabel = $derived(
    `Semana del ${formatDate(weekStart)} al ${formatDate(weekEnd)}`,
  );

  function toInputValue(date: Date): string {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
  }

  function formatDate(date: Date): string {
    const d = String(date.getDate()).padStart(2, "0");
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const y = date.getFullYear();
    return `${d}/${m}/${y}`;
  }

  function getWeekRange(date: Date): { start: Date; end: Date } {
    const d = new Date(date);
    const day = d.getDay();
    const diffToMonday = day === 0 ? -6 : 1 - day;
    const monday = new Date(d);
    monday.setDate(d.getDate() + diffToMonday);
    monday.setHours(0, 0, 0, 0);
    const sunday = new Date(monday);
    sunday.setDate(monday.getDate() + 6);
    sunday.setHours(0, 0, 0, 0);
    return { start: monday, end: sunday };
  }

  function handleDateChange(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    if (!target.value) return;
    const [y, m, d] = target.value.split("-").map(Number);
    const picked = new Date(y, m - 1, d);
    const range = getWeekRange(picked);
    setWeekRange(range.start, range.end);
  }
</script>

<div class="flex items-center gap-3">
  <Input
    type="date"
    class="w-40"
    value={inputValue}
    onchange={handleDateChange}
  />
  <span class="text-sm font-medium text-muted-foreground">{weekLabel}</span>
</div>
