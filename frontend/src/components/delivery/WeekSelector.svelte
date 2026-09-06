<script lang="ts">
  import DatePicker from "$components/ui/DatePicker.svelte";
  import {
    getSelectedDate,
    setSelectedDate,
  } from "$lib/stores/delivery.svelte.js";
  import { formatDate, displayToISO } from "$lib/utils/date.js";

  let dateValue = $state("");
  let syncing = false;

  $effect(() => {
    const selectedDate = getSelectedDate();
    syncing = true;
    dateValue = formatDate(selectedDate);
    syncing = false;
  });

  $effect(() => {
    if (!syncing && dateValue) {
      const iso = displayToISO(dateValue);
      const [y, m, d] = iso.split("-").map(Number);
      setSelectedDate(new Date(y, m - 1, d));
    }
  });
</script>

<div class="flex items-center gap-3">
  <span class="text-sm text-muted-foreground">Fecha:</span>
  <DatePicker bind:value={dateValue} class="w-40" />
</div>
