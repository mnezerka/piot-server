<script>
    import { createEventDispatcher } from 'svelte';

    export let value = null;

    const dispatch = createEventDispatcher();

    function onClick() {
        dispatch('click', {value});
    }

</script>

<a class="button is-danger" on:click={onClick}>
    <span class="icon is-small"><i class="fas fa-times"></i></span>
</a>
