# Fuzzing API Testing

Este proyecto realiza pruebas de fuzzing en una API REST utilizando Go, con el objetivo de descubrir posibles vulnerabilidades o fallos en el comportamiento de la API al recibir entradas inesperadas o malformadas.

Las pruebas de fuzzing se realizan tanto en los endpoints GET como POST de la API.

## Estructura del Proyecto

La estructura del proyecto es la siguiente:

fuzzing-api/ │ ├── api/ # Cliente de la API (para hacer peticiones GET y POST) ├── logger/ # Logger para registrar las semillas generadas durante las pruebas ├── fuzz/ # Pruebas de fuzzing para los endpoints de la API │ ├── fuzz_get_test.go # Test de fuzzing para la petición GET │ └── fuzz_post_test.go# Test de fuzzing para la petición POST ├── main.go # Punto de entrada del programa └── README.md # Este archivo

less
Copiar código

## Instalación

Para ejecutar este proyecto, necesitas tener Go instalado en tu máquina. Si aún no lo tienes, puedes descargarlo desde [aquí](https://golang.org/dl/).

1. **Clonar el repositorio**:

   Si aún no has clonado el repositorio, puedes hacerlo con el siguiente comando:

   ```bash
   git clone https://github.com/tu-usuario/fuzzing-api.git
   cd fuzzing-api
Instalar dependencias:

Asegúrate de que Go esté instalado en tu sistema y luego instala las dependencias necesarias:

bash
Copiar código
go mod tidy
Ejecución de las Pruebas de Fuzzing
Este proyecto incluye pruebas de fuzzing tanto para el endpoint GET como para el endpoint POST de la API.

1. Pruebas de Fuzzing para GET
Para ejecutar las pruebas de fuzzing para el endpoint GET, utiliza el siguiente comando:

bash
Copiar código
go test -fuzz=FuzzGetEndpoint -fuzztime=30s
Esto ejecutará las pruebas de fuzzing en el endpoint GET y generará un reporte HTML con las semillas generadas.

2. Pruebas de Fuzzing para POST
Para ejecutar las pruebas de fuzzing para el endpoint POST, utiliza el siguiente comando:

bash
Copiar código
go test -fuzz=FuzzPostEndpoint -fuzztime=30s
Al igual que para el GET, las pruebas para el endpoint POST también generarán un reporte HTML con las semillas utilizadas.

Reportes HTML
Después de ejecutar las pruebas de fuzzing, se generarán archivos HTML con los reportes. Los reportes contienen las semillas generadas durante las pruebas.

El reporte de las pruebas GET será guardado en un archivo llamado fuzz_get_report.html.
El reporte de las pruebas POST será guardado en un archivo llamado fuzz_post_report.html.
Puedes abrir estos archivos en cualquier navegador para ver las semillas generadas durante las pruebas.

Explicación de la Estructura de las Pruebas
Fuzzing para el Endpoint GET
Realiza peticiones al endpoint /Activities/{id}, donde {id} es una semilla que puede contener caracteres inesperados, como caracteres Unicode o secuencias malformadas.

Fuzzing para el Endpoint POST
Envía datos malformados a la API en el cuerpo de la solicitud POST. La API de ejemplo espera un JSON con un objeto que tiene campos como id, title, dueDate y completed.

Registro de Semillas
Durante las pruebas, todas las URL y datos de las solicitudes generadas se registran. Los detalles de cada semilla generada se incluyen en el reporte HTML generado al final de las pruebas.

Contribuciones
Si deseas contribuir al proyecto, siéntete libre de hacer un fork del repositorio y enviar pull requests.

Contacto
Si tienes alguna duda o pregunta, no dudes en abrir un issue en el repositorio o contactarnos directamente.

markdown
Copiar código

### Instrucciones para crear el archivo

1. **Crear el archivo README.md**: 
   Abre tu editor de texto preferido, crea un archivo nuevo llamado `README.md` y copia el contenido anterior en el archivo.

2. **Subir a tu repositorio**:
   Si ya tienes un repositorio en GitHub, puedes subir este archivo con los siguientes comandos:

   ```bash
   git add README.md
   git commit -m "Añadir README con instrucciones"
   git push origin main