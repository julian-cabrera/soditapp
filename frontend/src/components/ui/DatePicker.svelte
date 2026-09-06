<script lang="ts">
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Calendar } from "$lib/components/ui/calendar/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { parseDate, type DateValue } from "@internationalized/date";
  import { displayToISO, isoToDisplay } from "$lib/utils/date.js";

  let {
    value = $bindable(""),
    label,
    class: className,
  }: {
    value: string;
    label?: string;
    class?: string;
  } = $props();

  let dateValue = $derived.by(() => {
    if (!value) return undefined;
    try {
      return parseDate(displayToISO(value)) as DateValue;
    } catch {
      return undefined;
    }
  });

  function handleValueChange(date: DateValue | undefined) {
    if (date) {
      value = isoToDisplay(date.toString());
    }
  }
</script>

<div class="grid gap-1 {className}">
  {#if label}
    <Label>{label}</Label>
  {/if}
  <Popover.Root>
    <Popover.Trigger>
      <Button
        variant="outline"
        class="w-40 justify-start text-left font-normal"
      >
        {value || "Seleccionar fecha"}
      </Button>
    </Popover.Trigger>
    <Popover.Content class="w-auto p-0" align="start">
      <Calendar
        type="single"
        value={dateValue}
        onValueChange={handleValueChange}
      />
    </Popover.Content>
  </Popover.Root>
</div>
