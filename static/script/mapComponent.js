(function () {
  const ZOOM = 13;
  const CENTER_LATITUDE = 7.447220876004841;
  const CENTER_LONGITUDE = 125.80942522476553;

  map = new OpenLayers.Map("map");
  const mapnik = new OpenLayers.Layer.OSM();
  const fromProjection = new OpenLayers.Projection("EPSG:4326"); // Transform from WGS 1984
  const toProjection = new OpenLayers.Projection("EPSG:900913"); // to Spherical Mercator Projection
  const position = new OpenLayers.LonLat(
    CENTER_LONGITUDE,
    CENTER_LATITUDE,
  ).transform(fromProjection, toProjection);

  let markers = new OpenLayers.Layer.Markers("Markers");

  const coordinates = JSON.parse(
    document.getElementById("markers").textContent,
  );

  for (const coordinate of coordinates) {
    const marker = new OpenLayers.LonLat(
      coordinate.longitude,
      coordinate.latitude,
    ).transform(fromProjection, toProjection);
    markers.addMarker(new OpenLayers.Marker(marker));
  }

  map.addLayer(mapnik);
  map.addLayer(markers);
  map.setCenter(position, ZOOM);
})();
