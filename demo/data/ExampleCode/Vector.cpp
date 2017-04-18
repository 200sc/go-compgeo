// Vector.cpp: implementation of the Vector class.
//
//////////////////////////////////////////////////////////////////////

#include "stdafx.h"
#include "Vector.h"

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#define new DEBUG_NEW
#endif

/*
 * 3D Vector class. Encapsulates common 3D operations on vectors.
 * Author: Ryan Holmes
 * E-mail: ryan <at> holmes3d <dot> net
 */

//////////////////////////////////////////////////////////////////////
// Construction/Destruction
//////////////////////////////////////////////////////////////////////

Vector::Vector() : x(0), y(0), z(0)
{
}

Vector::Vector(double newX, double newY, double newZ) : x(newX), y(newY), z(newZ)
{
}

Vector::Vector(double array[]) : x(array[0]), y(array[1]), z(array[2])
{
}

Vector::Vector(float array[]) : x(array[0]), y(array[1]), z(array[2])
{
}

Vector::Vector(int array[]) : x(array[0]), y(array[1]), z(array[2])
{
}

Vector::~Vector()
{

}

#ifdef DIRECT3D_VERSION
Vector::Vector(D3DXVECTOR3 d3dV) : x(d3dV.x), y(d3dV.y), z(d3dV.z)
{
}

Vector::operator D3DXVECTOR3() const
{
	return D3DXVECTOR3((FLOAT)x,(FLOAT)y,(FLOAT)z);
}
#endif

Vector Vector::operator+(const Vector& rhs) const
{
    return Vector(x + rhs.x, y + rhs.y, z + rhs.z);
}

Vector Vector::operator-(const Vector& rhs) const
{
    return Vector(x - rhs.x, y - rhs.y, z - rhs.z);
}

Vector Vector::operator*(const double factor) const
{
    return Vector(x * factor, y * factor, z * factor);
}

Vector Vector::operator/(const double factor) const
{
    return Vector(x / factor, y / factor, z / factor);
}

double Vector::operator[](const int index) const
{
	switch (index) {
		case VECTOR_X:
			return x;
			break;
		case VECTOR_Y:
			return y;
			break;
		default:
			return z;
			break;
	}
}

double& Vector::operator[](const int index)
{
	switch (index) {
		case VECTOR_X:
			return x;
			break;
		case VECTOR_Y:
			return y;
			break;
		default:
			return z;
			break;
	}
}

void Vector::translateBy(const Vector& rhs)
{
	x += rhs.x;
	y += rhs.y;
	z += rhs.z;
}

void Vector::scaleBy(const double factor)
{
	x *= factor;
	y *= factor;
	z *= factor;
}

void Vector::normalize()
{
	double length = sqrt(x * x + y * y + z * z);
	if (length > 0.0001) {
		x /= length;
		y /= length;
		z /= length;
	}
}

Vector& Vector::normalized()
{
	double length = sqrt(x * x + y * y + z * z);
	if (length > 0.0001) {
		x /= length;
		y /= length;
		z /= length;
	}
	return *this;
}

double Vector::normalizeAndReturn()
{
	double length = sqrt(x * x + y * y + z * z);
	if (length > 0.0001) {
		x /= length;
		y /= length;
		z /= length;
	}
	return length;
}

void Vector::zero()
{
	x = y = z = 0.0;
}

double Vector::getLength() const
{
	return sqrt(x * x + y * y + z * z);
}

double Vector::getSquaredLength() const
{
	return (x * x + y * y + z * z);
}

double Vector::Dot(const Vector& rhs) const
{
	return (x * rhs.x + y * rhs.y + z * rhs.z);
}

Vector Vector::Cross(const Vector& rhs) const
{
	return Vector((y * rhs.z) - (z * rhs.y),
		          (z * rhs.x) - (x * rhs.z),
				  (x * rhs.y) - (y * rhs.x));
}

void Vector::toArray(double array[]) const
{
	array[0] = x;
	array[1] = y;
	array[2] = z;
}

void Vector::toArray(float array[]) const
{
	array[0] = (float)x;
	array[1] = (float)y;
	array[2] = (float)z;
}

void Vector::fromArray(double array[])
{
	x = array[0];
	y = array[1];
	z = array[2];
}

void Vector::fromArray(float array[])
{
	x = array[0];
	y = array[1];
	z = array[2];
}

void Vector::rotateX(const double degrees) {
	radianRotateX(degrees * DEG_TO_RAD);
}

void Vector::rotateY(const double degrees) {
	radianRotateY(degrees * DEG_TO_RAD);
}

void Vector::rotateZ(const double degrees) {
	radianRotateZ(degrees * DEG_TO_RAD);
}

void Vector::radianRotateX(const double radians) {
	double cosAngle = cos(radians);
	double sinAngle = sin(radians);
	double origY = y;
	y =	y * cosAngle - z * sinAngle;
	z = origY * sinAngle + z * cosAngle;
}

void Vector::radianRotateY(const double radians) {
	double cosAngle = cos(radians);
	double sinAngle = sin(radians);
	double origX = x;
	x =	x * cosAngle + z * sinAngle;
	z = z * cosAngle - origX * sinAngle;
}

void Vector::radianRotateZ(const double radians) {
	double cosAngle = cos(radians);
	double sinAngle = sin(radians);
	double origX = x;
	x =	x * cosAngle - y * sinAngle;
	y = origX * sinAngle + y * cosAngle;
}

void Vector::rotateAxis(const Vector& axis, const double degrees)
{
	radianRotateAxis(axis, degrees * DEG_TO_RAD);
}

void Vector::radianRotateAxis(const Vector& axis, const double radians)
{
	// Formula goes CW around axis. I prefer to think in terms of CCW
	// rotations, to be consistant with the other rotation methods.
	double cosAngle = cos(-radians);
	double sinAngle = sin(-radians);

	Vector w = axis;
	w.normalize();
	double vDotW = Dot(w);
	Vector vCrossW = Cross(w);
	w.scaleBy(vDotW); // w * (v . w)

	x = w.x + (x - w.x) * cosAngle + vCrossW.x * sinAngle;
	y = w.y + (y - w.y) * cosAngle + vCrossW.y * sinAngle;
	z = w.z + (z - w.z) * cosAngle + vCrossW.z * sinAngle;
}

Vector Vector::interpolate1(const Vector& endPoint, const double t) const
{
	return Vector(x + t * (endPoint.x - x),
		          y + t * (endPoint.y - y),
				  z + t * (endPoint.z - z));
}

Vector Vector::interpolate2(const Vector& midControl, const Vector& endControl, const double t) const
{
    Vector left = this->interpolate1(midControl, t);
	Vector right = midControl.interpolate1(endControl, t);
	return left.interpolate1(right, t);
}

Vector Vector::interpolate3(const Vector& leftControl, const Vector& rightControl, const Vector& endControl, const double t) const
{
    Vector begin = this->interpolate1(leftControl, t);
	Vector mid = leftControl.interpolate1(rightControl, t);
	Vector end = rightControl.interpolate1(endControl, t);
	return begin.interpolate2(mid, end, t);
}
