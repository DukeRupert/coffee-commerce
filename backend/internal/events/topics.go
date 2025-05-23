// internal/events/topics.go
package events

// Product-related topics
const (
	TopicProductCreated      = "products.created"
	TopicProductUpdated      = "products.updated"
	TopicProductDeleted      = "products.deleted"
	TopicProductStockUpdated = "products.stock_updated"
	TopicProductLowStock     = "products.low_stock"
)

// Product-related topics
const (
	TopicVariantCreated = "variants.created"
	TopicVariantQueued  = "variants.queued" // Event when a variant needs to be created
	TopicVariantUpdated = "variants.updated"
	TopicVariantDeleted = "variants.deleted"
)

// Customer-related topics
const (
	TopicCustomerCreated = "customers.created"
	TopicCustomerUpdated = "customers.updated"
	TopicCustomerDeleted = "customers.deleted"
)

// Subscription-related topics
const (
	TopicSubscriptionCreated  = "subscriptions.created"
	TopicSubscriptionUpdated  = "subscriptions.updated"
	TopicSubscriptionCanceled = "subscriptions.canceled"
	TopicSubscriptionPaused   = "subscriptions.paused"
	TopicSubscriptionResumed  = "subscriptions.resumed"
	TopicSubscriptionRenewed  = "subscriptions.renewed"
)

// Order-related topics
const (
	TopicOrderCreated       = "orders.created"
	TopicOrderStatusUpdated = "orders.status_updated"
	TopicOrderShipped       = "orders.shipped"
	TopicOrderDelivered     = "orders.delivered"
)

// Stripe-related topics
const (
	TopicStripeProductCreated = "stripe.products.created"
	TopicStripeProductUpdated = "stripe.products.updated"
	TopicStripeProductDeleted = "stripe.products.deleted"
	
	TopicStripePriceCreated = "stripe.prices.created"
	TopicStripePriceUpdated = "stripe.prices.updated"
	TopicStripePriceDeleted = "stripe.prices.deleted"
	
	TopicStripeCustomerCreated = "stripe.customers.created"
	TopicStripeCustomerUpdated = "stripe.customers.updated"
	TopicStripeCustomerDeleted = "stripe.customers.deleted"
	
	TopicStripeSubscriptionCreated = "stripe.subscriptions.created"
	TopicStripeSubscriptionUpdated = "stripe.subscriptions.updated"
	TopicStripeSubscriptionCanceled = "stripe.subscriptions.canceled"
	
	TopicStripeCheckoutCompleted = "stripe.checkout.completed"
	TopicStripeInvoicePaid = "stripe.invoice.paid"
	TopicStripeInvoicePaymentFailed = "stripe.invoice.payment_failed"
)
