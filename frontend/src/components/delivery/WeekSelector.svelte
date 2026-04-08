<script lang="ts">
  import { Input } from "$lib/components/ui/input/index.js";
  import {
    getSelectedDate,
    setSelectedDate,
  } from "$lib/stores/delivery.svelte.js";

  let selectedDate = $derived(getSelectedDate());

  let inputValue = $derived(toInputValue(selectedDate));

  function toInputValue(date: Date): string {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
  }

  function handleDateChange(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    if (!target.value) return;
    const [y, m, d] = target.value.split("-").map(Number);
    setSelectedDate(new Date(y, m - 1, d));
  }
</script>

<div class="flex items-center gap-3">
  <span class="text-sm text-muted-foreground">Fecha:</span>
  <Input
    type="date"
    class="w-40"
    value={inputValue}
    onchange={handleDateChange}
  />
</div>
