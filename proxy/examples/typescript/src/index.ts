// proxy examples in TypeScript
import { dag, object, func, Service, Proxy } from "@dagger.io/dagger";

@object()
class Example {
  /**
   * example for withservice function
   */
  @func()
  proxyWithService(service: Service): Service {
    return dag.proxy().withService(service, "myService", 8080, 80).service();
  }

  /**
   * example for service function
   */
  @func()
  proxyService(serviceA: Service, serviceB: Service): Service {
    return dag
      .proxy()
      .withService(serviceA, "serviceA", 8080, 80)
      .withService(serviceB, "serviceB", 8081, 80)
      .service();
  }
}
