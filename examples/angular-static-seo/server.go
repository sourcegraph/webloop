package main

import (
	"flag"
	"github.com/sourcegraph/webloop"
	"github.com/sqs/gotk3/gtk"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var appBind = flag.String("app-http", ":9000", "HTTP bind address for AngularJS app")
var staticBind = flag.String("static-http", ":9100", "HTTP bind address for static app")

func main() {
	flag.Parse()

	appMux := http.NewServeMux()
	appMux.HandleFunc("/", serveApp)
	go start("app", *appBind, appMux)

	staticMux := http.NewServeMux()
	staticMux.HandleFunc("/", serveStatic)
	start("static", *staticBind, staticMux)
}

func start(name, bind string, mux *http.ServeMux) {
	log.Printf("%s: Listening on %s", name, bind)
	err := http.ListenAndServe(bind, mux)
	if err != nil {
		log.Fatalf("%s: ListenAndServe: %s", name, err)
	}
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	w.Write(page)
}

func init() {
	gtk.Init(nil)
	go func() {
		runtime.LockOSThread()
		gtk.Main()
	}()
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	var ctx webloop.Context

	view := ctx.NewView()
	defer view.Close()

	r.URL.Host = "localhost" + *appBind
	r.URL.Scheme = "http"
	log.Printf("Generating static page for URL: %s", r.URL)
	view.Open(r.URL.String())
	view.Wait()

	// Wait until window.$viewReadyForSnapshot is true.
	timeout := time.Second * 3
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			http.Error(w, "application did not set $viewReadyForSnapshot within timeout "+timeout.String(), http.StatusInternalServerError)
			return
		}

		check, err := view.EvaluateJavaScript("window.$viewReadyForSnapshot")
		if err != nil {
			http.Error(w, "error checking $viewReadyForSnapshot: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ready, _ := check.GoValue()
		if ready, ok := ready.(bool); !ok || !ready {
			time.Sleep(timeout / 30)
			continue
		}

		result, err := view.EvaluateJavaScript("document.documentElement.outerHTML")
		if err != nil {
			http.Error(w, "error generating static page: "+err.Error(), http.StatusInternalServerError)
			return
		}
		html := result.String()
		html = strings.Replace(html, "<body>", `<body><h3>This is a static page generated from <a href="`+r.URL.String()+`">`+r.URL.String()+`</a></h3><hr>`, 1)
		html = strings.Replace(html, "ng-app=", "disabled-ng-app=", -1)
		html = strings.Replace(html, "</pre>", "\nGenerated static page in "+time.Since(start).String()+"\n</pre>", 1)
		w.Write([]byte(html))
		return
	}

}

var page = []byte(`
<!doctype html>
<html ng-app="staticSEO">
<head>
  <meta charset="utf-8">
  <title>WebLoop angular-static-seo example</title>
  <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.2.0-rc.3/angular.min.js"></script>
  <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.2.0-rc.3/angular-route.min.js"></script>
</head>
<body>

<div ng-view></div>

<hr>

<p><a style="color: #777" href="https://sourcegraph.com/github.com/sourcegraph/webloop/readme">WebLoop example: AngularJS static SEO</a></p>

<pre>
URL:         {{$location.url()}}

Params:      {{$route.current.params}}

User-Agent:  {{userAgent}}
</pre>

<script type=text/ng-template id="index.html">
  <h2>Angular static SEO example</h2>
  <p>
    This sample <a href="http://angularjs.org">AngularJS</a> application demonstrates how to use
    <a href="https://sourcegraph.com/github.com/sourcegraph/webloop/readme">WebLoop</a> to
    generate a static, SEO-friendly site from a single-page AngularJS application.
  </p>
  <hr>
  <h1>Cities</h1>
  <p>Showing {{cities.length}} cities.</p>
  <ul>
    <li ng-repeat="city in cities">
      <a ng-href="/cities/{{city.id}}">{{city.name}}</a> (population: {{city.population}})
    </li>
  </ul>
</script>

<script type=text/ng-template id="detail.html">
  <p><a href="/cities">&laquo; Back to list of cities</a></p>
  <h1>{{city.name}}</h1>
  <table>
    <tr><th>Population:</th><td>{{city.population}}</td></tr>
  </table>
</script>

<script>
var allCities = [
  {id: 'shanghai', name: 'Shanghai', population: 17836133},
  {id: 'istanbul', name: 'Istanbul', population: 13854740},
  {id: 'karachi', name: 'Karachi', population: 12991000},
  {id: 'mumbai', name: 'Mumbai', population: 12478447},
  {id: 'moscow', name: 'Moscow', population: 11977988},
  {id: 'sao-paulo', name: 'Sao Paulo', population: 11821876},
  {id: 'beijing', name: 'Beijing', population: 11716000},
];

angular.module('staticSEO', ['ngRoute'])

.config(function($locationProvider, $routeProvider) {
  $locationProvider.html5Mode(true);

  $routeProvider
    .when('/cities/:city', {
      controller: 'CityCtrl',
      resolve: {
        city: function($q, $route, $timeout) {
          var cityID = $route.current.params.city;
          var deferred = $q.defer();
          // Simulate loading delay.
          $timeout(function() {
            var city = allCities.filter(function(city) {
              return city.id === cityID;
            })[0];
            if (city) deferred.resolve(city); 
            else deferred.reject('No city found with ID "' + cityID + '"');
          }, 500);
          return deferred.promise;
        },
      },
      templateUrl: 'detail.html',
    })
    .when('/cities', {
      controller: 'CitiesCtrl',
      templateUrl: 'index.html',
    })
    .otherwise({
      redirectTo: '/cities',
    });
})

.run(function($location, $rootScope, $route, $window) {
  $rootScope.userAgent = $window.navigator.userAgent;
  $rootScope.$location = $location;
  $rootScope.$route = $route;

  $rootScope.$on('$viewReadyForSnapshot', function() {
    $window.$viewReadyForSnapshot = true;
  });
  $rootScope.$on('$routeChangeBegin', function() {
    $window.$viewReadyForSnapshot = false;
  });
})

.controller('CitiesCtrl', function($scope, $timeout) {
  $timeout(function() {
    $scope.cities = allCities;
    $scope.$emit('$viewReadyForSnapshot');
  }, 350);
})

.controller('CityCtrl', function($scope, city) {
  $scope.city = city;
  $scope.$emit('$viewReadyForSnapshot');
})

;
</script>
</body>
</html>
`)
