<script>
    import {onMount} from 'svelte';
    import {authenticated, token} from '../stores'
    import {push} from 'svelte-spa-router'

    onMount(() => {
        if (!$authenticated) { push("/login"); }
    })
</script>

<h1 class="title">Things</h1>

