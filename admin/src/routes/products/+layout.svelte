<script lang="ts">
  import { page } from '$app/state';
</script>

<div class="products-layout">
  <header class="products-header">
    <div class="container">
      <h1>Coffee Products</h1>
      
      <nav class="products-nav">
        <ul>
          <li class:active={page.url.pathname === '/products'}>
            <a href="/products">All Products</a>
          </li>
          <li class:active={page.url.pathname === '/products/new'}>
            <a href="/products/create">Create New</a>
          </li>
        </ul>
      </nav>
    </div>
  </header>
  
  <main>
    <!-- Success message shown after redirect with ?created=true -->
    {#if page.url.searchParams.get('created') === 'true'}
      <div class="container">
        <div class="alert success">
          <p>Product created successfully!</p>
        </div>
      </div>
    {/if}
    
    <slot />
  </main>
</div>

<style>
  .products-header {
    background-color: #f9fafb;
    border-bottom: 1px solid #e5e7eb;
    padding: 1.5rem 0;
    margin-bottom: 2rem;
  }
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
  }
  
  .products-header .container {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  @media (min-width: 768px) {
    .products-header .container {
      flex-direction: row;
      justify-content: space-between;
      align-items: center;
    }
  }
  
  .products-header h1 {
    margin: 0;
    color: #111827;
    font-size: 1.5rem;
  }
  
  .products-nav ul {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 1.5rem;
  }
  
  .products-nav a {
    text-decoration: none;
    color: #6b7280;
    font-weight: 500;
    padding: 0.5rem 0;
    position: relative;
  }
  
  .products-nav a:hover {
    color: #4b5563;
  }
  
  .products-nav li.active a {
    color: #3b82f6;
  }
  
  .products-nav li.active a::after {
    content: '';
    position: absolute;
    bottom: -2px;
    left: 0;
    right: 0;
    height: 2px;
    background-color: #3b82f6;
    border-radius: 2px;
  }
  
  .alert {
    margin-bottom: 2rem;
    padding: 1rem;
    border-radius: 4px;
  }
  
  .alert.success {
    background-color: #dcfce7;
    color: #166534;
    border: 1px solid #bbf7d0;
  }
</style>