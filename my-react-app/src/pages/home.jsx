import Navbar from "../components/Navbar.jsx";

import ProductCard from "../components/ProductCard.jsx";

function Home() {


  return (
    <>
      <div>
        <Navbar></Navbar>
      </div>
      <div className="p-2">
        <ProductCard></ProductCard>
      </div>
    </>
  );
}

export default Home;
