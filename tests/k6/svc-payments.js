import http from "k6/http";
import { group, check } from "k6";
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.6.0/index.js';
import { listSubscriptions } from "./svc-subscriptions.js";

const SVC_PATH = "/payments";
const BASE_URL = `${__ENV.PAYMENTS_ENDPOINT}`;

const getURL = () => {
  return fixEndpoint(BASE_URL, SVC_PATH);
}

function fixEndpoint(endpoint, suffixPath) {
  if (!endpoint.startsWith("http")) {
    if (endpoint.startsWith(":")) {
      return `http://localhost${endpoint}${!endpoint.endsWith(suffixPath) ? suffixPath : ""}`; 
    } else {
      return `http://localhost:8080/${endpoint}${!endpoint.endsWith(suffixPath) ? suffixPath : ""}`;
    }
  }
  return endpoint;
}

const listPayments = () => {
  const res = http.get(getURL(), { headers: { "Content-Type": "application/json", "Accept": "application/json" } });
  return res;
};

const paymentStatus = () => Math.random() < 0.5 ? "SUCCESS" : "FAILED";

const generatePayment = (subscriptionId, amount) => ({
  id: uuidv4(),
  subscription_id: subscriptionId,
  amount,
  status: paymentStatus(),
});

const makePayment = (subscriptionId, amount) => {
  const payload = JSON.stringify(generatePayment(subscriptionId, amount));
  const res = http.post(getURL(), payload, { headers: { "Content-Type": "application/json" } });
  return res;
};

export function testListPayments() {
  group("list payments", () => {
    const res = listPayments();
    check(res, { "Successful listing": (r) => r.status === 200 });
  });
}

export function testMakePayments(qty) {
  group("make payments", () => {
    for (let i = 0; i < qty; i++) {
      const response = listSubscriptions();
      if (response.status !== 200) {
        console.error("Failed to list subscriptions:", response.error_code);
        continue;
      }

      try {
        const subscriptions = response.json();
        const randomSubscription = subscriptions[Math.floor(Math.random() * subscriptions.length)];
        const subscriptionId = randomSubscription.id;

        const amount = Math.floor(Math.random() * 10);

        const paymentResponse = makePayment(subscriptionId, amount);

        check(paymentResponse, {  // Corrected status code check (201)
          "Payment created": (r) => r.status === 201,
        });

        if (paymentResponse.status !== 201) {
          console.error("Failed to make a payment:", paymentResponse.error_code); // Include error for failed payment creation
        }

      } catch (error) {
        console.error("Error processing subscriptions:", error);
      }
    }
  });
}


