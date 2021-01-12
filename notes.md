# Notes

- Cost Explorer only has two REGION endpoints, us-east-1 and cn-northwest-1. See [this](https://docs.aws.amazon.com/general/latest/gr/billing.html)
- Pricing can potentially get out of hand:

1 Request = 0.01 USD
Current amount of Capability contexts = 118

Assuming each context/AWS account was queried once every hour for 24 hours:
1 hour: (118 * 0.01) = 1.18 USD
24 hours: 1.18 * 24 = 28.32 USD

Now, if that was done once every hour for 30 days:
30 days: 28.32 * 30 = 847 USD

Seems like a lot of money for getting to know how much money we spend.
