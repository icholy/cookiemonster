<?php
  // wafer hook for PHP. Simply set the JWT cookie
  if (isset($_GET["jwt"])) {
    $jwt = $_GET["jwt"];
    setcookie("jwt", $jwt);
  }
?>
