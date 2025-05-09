// Testes básicos acessando rotas dos serviços
// VU = virtual user
import { testAddUsers} from "./svc-users.js";
import { testAddPlans} from "./svc-plans.js";
import { testAddSubscriptions} from "./svc-subscriptions.js";
import { testMakePayments} from "./svc-payments.js";

// https://grafana.com/docs/k6/latest/using-k6/k6-options/reference/
export const options = {
    // Key configurations for spike in this section
    // stages: [
    //   // A list of VU { target: ..., duration: ... } objects that specify the target number of VUs to ramp up or down to for a specific period.
    //   { duration: '2m', target: 20 },
    //   { duration: '1m', target: 1 },
    //   { duration: '2m', target: 20 },
    // ],
    // noConnectionReuse: true,
    insecureSkipTLSVerify: true,
    userAgent: 'k6-otel/1.0',
  
    httpDebug: '',
  
    summaryTimeUnit: 's',

    iterations: 200,
    duration: '1m',
};

// Global variables should be initialized.

export default function() {

    testAddUsers(20);

    testAddPlans(20);

    testAddSubscriptions(20);

    testMakePayments(10);
}


