// Params

numberOfTeeth=10;
teethHeightEqualDivideTeethWidthBy = 2.5;
extrusion=1;

// Constants

radius=100;
firstAngle=-45;
$fn=64;

// Calculated

wheelCirc = 2*PI*radius;
echo("wheel circ: ", wheelCirc);

teethWidth = wheelCirc/(numberOfTeeth*2);
echo("teeth width: ", teethWidth);
angleBetweehTeeth=360/numberOfTeeth;
echo("angle between theeth: ", angleBetweehTeeth);

wheelWidth=teethWidth;
echo("wheel width: ", wheelWidth);

teethHeight=teethWidth/teethHeightEqualDivideTeethWidthBy;
echo("teeth height: ", teethHeight);


// -  - .    ×
// |  B .  ×   ×
// |  - .×       ×
// |  | .  ×       ×
// R  | .    ×       ×
// |  A .    C ×       ×
// |  | .        ×   ×
// |  | .          ×   W
// -  - .................

// A²+A² = C²
// C² = 2×A²
// C = sqrt(2×A²)

// A = R - B

// R = radius + teethHeight;

// B² + B² = W²
// 2×B² = W²
// B² = W²/2
// B = sqrt(W²/2)

// W = teethWidth/2, because teethWidth goes both sides of the radius

www = teethWidth/2;

bbb = sqrt((www^2)/2);
rrr = radius + teethHeight;
aaa = rrr - bbb;
ccc = sqrt(2*(aaa^2));

echo("length to go in the angle: ", ccc);

footLength = 2*(ccc - (radius - (wheelWidth/2)));

// Drawing

module teeth(angle) {
    rotate(angle)
        translate([(radius+teethHeight)/2, 0, 0])
            square([radius+teethHeight, teethWidth], true);
    newangle = angle + angleBetweehTeeth;
    if (newangle < firstAngle+360) {
        teeth(newangle);
    }
}

color("mediumpurple") {
    difference() {
        union() {
            circle(radius);
            teeth(firstAngle);
        }
        
        circle(radius-wheelWidth);
    }

    rotate(firstAngle)
        translate([radius-wheelWidth/2, 0, 0])
            square([footLength, teethWidth], true);
}